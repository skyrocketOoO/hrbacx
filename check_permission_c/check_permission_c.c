#include "postgres.h"
#include "fmgr.h"
#include "utils/hsearch.h"
#include "access/htup_details.h"
#include "executor/spi.h"
#include "utils/array.h"

PG_MODULE_MAGIC;

typedef struct Queue {
    char **items;
    int front, rear, size, capacity;
} Queue;

// Function to create a queue
Queue* createQueue(int capacity) {
    Queue *queue = (Queue*) palloc(sizeof(Queue));
    queue->capacity = capacity;
    queue->size = 0;
    queue->front = 0;
    queue->rear = capacity - 1;
    queue->items = (char**) palloc(sizeof(char*) * capacity);
    return queue;
}

// Queue utility functions
int isQueueFull(Queue *queue) { return (queue->size == queue->capacity); }
int isQueueEmpty(Queue *queue) { return (queue->size == 0); }
void enqueue(Queue *queue, char *item) {
    if (isQueueFull(queue)) return;
    queue->rear = (queue->rear + 1) % queue->capacity;
    queue->items[queue->rear] = item;
    queue->size++;
}
char* dequeue(Queue *queue) {
    if (isQueueEmpty(queue)) return NULL;
    char *item = queue->items[queue->front];
    queue->front = (queue->front + 1) % queue->capacity;
    queue->size--;
    return item;
}

// Function to check permission
PG_FUNCTION_INFO_V1(check_permission_c);

Datum check_permission_c(PG_FUNCTION_ARGS) {
    text *user_id = PG_GETARG_TEXT_P(0);
    text *permission_type = PG_GETARG_TEXT_P(1);
    text *object_id = PG_GETARG_TEXT_P(2);

    // Convert text inputs to C strings
    char *user_id_str = text_to_cstring(user_id);
    char *permission_type_str = text_to_cstring(permission_type);
    char *object_id_str = text_to_cstring(object_id);

    // Initialize SPI (Server Programming Interface)
    if (SPI_connect() != SPI_OK_CONNECT) {
        elog(ERROR, "Failed to connect to SPI");
    }

    // Create a queue for BFS
    Queue *queue = createQueue(1000);  // Assuming a maximum of 1000 items
    enqueue(queue, user_id_str);

    // Set up a hash table for visited nodes
    HTAB *visited;
    HASHCTL hash_ctl;
    MemSet(&hash_ctl, 0, sizeof(hash_ctl));
    hash_ctl.keysize = sizeof(char*);
    hash_ctl.entrysize = sizeof(char*);
    visited = hash_create("Visited nodes", 1000, &hash_ctl, HASH_ELEM | HASH_FUNCTION);

    bool permission_found = false;

    while (!isQueueEmpty(queue)) {
        char *current = dequeue(queue);

        // Check if the current node has already been visited
        if (hash_search(visited, &current, HASH_FIND, NULL) != NULL) {
            continue;  // Skip if visited
        }

        // Mark current as visited
        hash_search(visited, &current, HASH_ENTER, NULL);

        // Query to check for permission
        char query[1024];
        snprintf(query, sizeof(query),
                 "SELECT 1 FROM \"edges\" WHERE from_v = '%s' AND relation = '%s' AND to_v = '%s'", 
                 current, permission_type_str, object_id_str);

        if (SPI_exec(query, 0) > 0 && SPI_processed > 0) {
            permission_found = true;
            break;
        }

        // Enqueue neighbors (leader_of relation)
        snprintf(query, sizeof(query),
                 "SELECT to_v FROM \"edges\" WHERE from_v = '%s' AND relation = 'leader_of'", current);

        if (SPI_exec(query, 0) > 0) {
            for (int i = 0; i < SPI_processed; i++) {
                char *to_v = SPI_getvalue(SPI_tuptable->vals[i], SPI_tuptable->tupdesc, 1);
                enqueue(queue, to_v);
            }
        }
    }

    // Clean up
    SPI_finish();
    if (permission_found) {
        PG_RETURN_BOOL(true);
    } else {
        PG_RETURN_BOOL(false);
    }
}
