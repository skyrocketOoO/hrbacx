import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = 'http://localhost:8080';
const headers = {
    'Content-Type': 'application/json',
}

export function AddLeader(leaderID, roleID){
    let response = http.post(`${BASE_URL}/addLeader`, JSON.stringify({
        leaderID: leaderID,
        roleID: roleID
    }), {
        headers: headers,
    });
    check(res, { 'AddLeader': (r) => r.status == 200 });
    return response;
}

export function AssignPermission(objectID, permissionType, roleID){
    let response = http.post(`${BASE_URL}/assignPermission`, JSON.stringify({
        objectID: objectID,
        permissionType: permissionType,
        roleID: roleID
    }), {
        headers: headers,
    });
    check(res, { 'AssignPermission': (r) => r.status == 200 });
    return response;
}

export function AssignRole(userID, roleID){
    let response = http.post(`${BASE_URL}/assignRole`, JSON.stringify({
        userID: userID,
        roleID: roleID
    }), {
        headers: headers,
    });
    check(res, { 'AssignRole': (r) => r.status == 200 });
    return response;
}

export function CheckPermission(userID, permissionType, objectID){
    let response = http.post(`${BASE_URL}/checkPermission`, JSON.stringify({
        userID: userID,
        permissionType: permissionType,
        objectID: objectID
    }), {
        headers: headers,
    });
    check(res, { 'CheckPermission': (r) => r.status == 200 });
    return response;
}