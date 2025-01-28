import { sleep, performance } from 'k6';
import { AddLeader, AssignPermission, AssignRole, CheckPermission, ClearAll } from './client.js';

export const options = {
  vus: 1,
  setupTimeout: '6000s',
}

const ROLE_CHILDREN = 5;
const ROLE_LAYER = 5;
const OBJECT_CHILDREN = 10;
const PERMISSION_TYPE = "write";

const sourceUser =  0;

export function setup() {
  const start = Date.now();
  console.log('Setting up the test...');

  AssignRole(sourceUser, 0);

  let cur = 0;
  for (let layer = 0; layer < ROLE_LAYER; layer++){
    let p = Math.pow(ROLE_CHILDREN, layer)
    for (let r = cur; r < cur+p; r++){
      for (let o = r * OBJECT_CHILDREN; o < (r+1) * OBJECT_CHILDREN; o++){
        AssignPermission(o, PERMISSION_TYPE, r);
      }
      if (layer == ROLE_LAYER-1){
        continue;
      }
      for (let c = r * ROLE_CHILDREN+1; c <= (r+1) * ROLE_CHILDREN; c++){
        AddLeader(r, c);
      }
    }
    cur += p;
  }

  const lastRole = cur-1;

  const setupTime = Date.now() - start;
  console.log(`Setup time: ${setupTime} ms`);
  sleep(5);
  return { lastRole };
}


export default function (data) {
  const { lastRole } = data;
  console.log(`last role: ${lastRole}`);

  const startTime = Date.now(); // Start measuring time

  let totalChecks = 0;
  for (let tar = 0; tar < lastRole * OBJECT_CHILDREN - 1; tar++) {
    CheckPermission(0, PERMISSION_TYPE, tar);
    totalChecks++;
  }

  const endTime = Date.now(); // End measuring time
  const totalTime = endTime - startTime; // Total time in milliseconds
  const averageTime = totalTime / totalChecks; // Average time per check

  console.log(`Total CheckPermission time: ${totalTime.toFixed(2)} ms`);
  console.log(`Average CheckPermission time: ${averageTime.toFixed(2)} ms`);
}


export function teardown(data) {
  ClearAll();
}