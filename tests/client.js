import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = 'http://localhost:8080';
const headers = {
    'Content-Type': 'application/json',
}

export function AddLeader(leaderID, roleID){
    let res = http.post(`${BASE_URL}/addLeader`, JSON.stringify({
        leaderID: leaderID.toString(),
        roleID: roleID.toString()
    }), {
        headers: headers,
    });
    check(res, { 'AddLeader': (r) => r.status == 200 });
    return res;
}

export function AssignPermission(objectID, permissionType, roleID){
    let res = http.post(`${BASE_URL}/assignPermission`, JSON.stringify({
        objectID: objectID.toString(),
        permissionType: permissionType,
        roleID: roleID.toString()
    }), {
        headers: headers,
    });
    check(res, { 'AssignPermission': (r) => r.status == 200 });
    return res;
}

export function AssignRole(userID, roleID){
    let res = http.post(`${BASE_URL}/assignRole`, JSON.stringify({
        userID: userID.toString(),
        roleID: roleID.toString()
    }), {
        headers: headers,
    });
    check(res, { 'AssignRole': (r) => r.status == 200 });
    return res;
}

export function CheckPermission(userID, permissionType, objectID){
    let res = http.post(`${BASE_URL}/checkPermission`, JSON.stringify({
        userID: userID.toString(),
        permissionType: permissionType,
        objectID: objectID.toString()
    }), {
        headers: headers,
    });
    check(res, { 'CheckPermission': (r) => r.status == 200 });
    return res;
}