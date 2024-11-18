import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

// Test configuration
export const options = {
    stages: [
        { duration: '1m', target: 500 },  // Ramp up to 500 users
        { duration: '3m', target: 500 },  // Stay at 500 users for 3 minutes
        { duration: '1m', target: 0 },    // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
        'http_req_duration{type:create}': ['p(95)<600'],
        'http_req_duration{type:get}': ['p(95)<400'],
        errors: ['rate<0.1'],             // Error rate should be less than 10%
    },
};

const BASE_URL = 'http://localhost:8080/api';

// Helper function to generate random user data
function generateRandomUser() {
    const randomId = Math.floor(Math.random() * 1000000);
    return {
        name: `Test User ${randomId}`,
        email: `test${randomId}@example.com`,
        password: `Test@${randomId}pass` // Updated to meet password requirements
    };
}

export default function () {
    const user = generateRandomUser();

    // Create user
    const createRes = http.post(`${BASE_URL}/users`, JSON.stringify(user), {
        headers: { 'Content-Type': 'application/json' },
        tags: { type: 'create' },
    });
    
    check(createRes, {
        'create user status is 201': (r) => r.status === 201,
    }) || errorRate.add(1);

    if (createRes.status === 201) {
        const userId = createRes.json('id');

        // Get user
        const getRes = http.get(`${BASE_URL}/users/${userId}`, {
            tags: { type: 'get' },
        });
        
        check(getRes, {
            'get user status is 200': (r) => r.status === 200,
            'get user returns correct email': (r) => r.json('email') === user.email,
        }) || errorRate.add(1);

        // Update user
        const updateRes = http.put(`${BASE_URL}/users/${userId}`, 
            JSON.stringify({ ...user, name: `Updated ${user.name}` }), {
            headers: { 'Content-Type': 'application/json' },
            tags: { type: 'update' },
        });
        
        check(updateRes, {
            'update user status is 200': (r) => r.status === 200,
        }) || errorRate.add(1);
    }

    // List users
    const listRes = http.get(`${BASE_URL}/users`, {
        tags: { type: 'list' },
    });
    
    check(listRes, {
        'list users status is 200': (r) => r.status === 200,
    }) || errorRate.add(1);

    // Small sleep to prevent overwhelming the server
    sleep(1);
}
