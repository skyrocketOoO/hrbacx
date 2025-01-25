import { sleep } from 'k6';
import { AddLeader, AssignPermission, AssignRole, CheckPermission, ClearAll } from './client.js';

export const options = {
  vus: 1,
  setupTimeout: '6000s',
}

const ROLE_CHILDREN = 6;
const ROLE_LAYER = 6;
const OBJECT_CHILDREN = 10;
const PERMISSION_TYPE = "write";

const sourceUser =  0;

export function setup() {
  const start = Date.now();
  console.log('Setting up the test...');

  AssignRole(sourceUser, 0);

  let cur = 0;
  for (let l = 0; l < ROLE_LAYER; l++){
    let p = Math.pow(ROLE_CHILDREN, l)
    for (let r = cur; r < cur+p; r++){
      for (let o = r * OBJECT_CHILDREN; o < (r+1) * OBJECT_CHILDREN; o++){
        // console.log(o, r)
        AssignPermission(o, PERMISSION_TYPE, r);
      }
      if (l == ROLE_LAYER-1){
        continue;
      }
      for (let c = r * ROLE_CHILDREN+1; c <= (r+1) * ROLE_CHILDREN; c++){
        AddLeader(r, c);
      }
    }
    cur += p;
  }

  const lastRole = cur-1;
  const targetObject = lastRole * OBJECT_CHILDREN-1;

  const setupTime = Date.now() - start;
  console.log(`Setup time: ${setupTime} ms`);
  sleep(5);
  return { targetObject };
}


export default function (data) {
  const { targetObject } = data; 
  console.log(`Target object: ${targetObject}`);
  CheckPermission(0, PERMISSION_TYPE,targetObject )
}

export function teardown(data) {
  ClearAll();
}