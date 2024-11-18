import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.2/index.js';

const errorRate = new Rate('errors');

// Custom metrics
const createUserTrend = new Trend('create_user_duration');
const getUserTrend = new Trend('get_user_duration');
const updateUserTrend = new Trend('update_user_duration');
const listUsersTrend = new Trend('list_users_duration');
const successfulRequests = new Counter('successful_requests');
const failedRequests = new Counter('failed_requests');

// Test configuration
export const options = {
    stages: [
        { duration: '1m', target: 500 },  // Ramp up to 500 users
        { duration: '3m', target: 500 },  // Stay at 500 users for 3 minutes
        { duration: '1m', target: 0 },  // Stay at 500 users for 3 minutes
    ],
    thresholds: {
        http_req_duration: ['p(95)<100'], // 95% of requests should be below 500ms
        'http_req_duration{type:create}': ['p(95)<200'],
        'http_req_duration{type:get}': ['p(95)<100'],
        errors: ['rate<0.1'],             // Error rate should be less than 10%
    },
};

const BASE_URL = 'http://localhost:8080/api';

// HTTP request configuration
const requestConfig = {
    timeout: '30s',
    headers: { 'Content-Type': 'application/json' },
};

// Retry configuration
const maxRetries = 3;
const retryDelay = 1000; // 1 second

// Helper function to generate random user data
function generateRandomUser() {
    const timestamp = Date.now();
    const random = Math.floor(Math.random() * 1000000);
    const uniqueId = `${timestamp}_${random}`;
    return {
        name: `Test User ${uniqueId}`,
        email: `test.${uniqueId}@example.com`,
        password: `Test@${uniqueId}123` // Meets all password requirements: uppercase, lowercase, number, special char
    };
}

// Helper function to perform request with retries
function performRequestWithRetry(fn) {
    let retries = 0;
    while (retries < maxRetries) {
        try {
            const response = fn();
            if (response.status !== 0) { // status 0 indicates network error
                return response;
            }
        } catch (error) {
            console.error(`Request failed (attempt ${retries + 1}/${maxRetries}): ${error}`);
        }
        retries++;
        if (retries < maxRetries) {
            sleep(retryDelay / 1000); // k6 sleep expects seconds
        }
    }
    return fn(); // last attempt
}

export default function () {
    const user = generateRandomUser();

    // Create user with retry
    const createRes = performRequestWithRetry(() => 
        http.post(`${BASE_URL}/users`, JSON.stringify(user), {
            ...requestConfig,
            tags: { type: 'create' },
        })
    );
    
    createUserTrend.add(createRes.timings.duration);
    
    check(createRes, {
        'create user status is 201': (r) => {
            const isSuccess = r.status === 201;
            isSuccess ? successfulRequests.add(1) : failedRequests.add(1);
            return isSuccess;
        },
    }) || errorRate.add(1);

    if (createRes.status === 201) {
        const userId = createRes.json('id');

        // Get user with retry
        const getRes = performRequestWithRetry(() =>
            http.get(`${BASE_URL}/users/${userId}`, {
                ...requestConfig,
                tags: { type: 'get' },
            })
        );
        
        getUserTrend.add(getRes.timings.duration);
        
        check(getRes, {
            'get user status is 200': (r) => {
                const isSuccess = r.status === 200;
                isSuccess ? successfulRequests.add(1) : failedRequests.add(1);
                return isSuccess;
            },
            'get user returns correct email': (r) => r.json('email') === user.email,
        }) || errorRate.add(1);

        // Update user with retry
        const updateRes = performRequestWithRetry(() =>
            http.put(`${BASE_URL}/users/${userId}`, 
                JSON.stringify({ ...user, name: `Updated ${user.name}` }), {
                ...requestConfig,
                tags: { type: 'update' },
            })
        );
        
        updateUserTrend.add(updateRes.timings.duration);
        
        check(updateRes, {
            'update user status is 200': (r) => {
                const isSuccess = r.status === 200;
                isSuccess ? successfulRequests.add(1) : failedRequests.add(1);
                return isSuccess;
            },
        }) || errorRate.add(1);
    }

    // List users with retry
    const listRes = performRequestWithRetry(() =>
        http.get(`${BASE_URL}/users`, {
            ...requestConfig,
            tags: { type: 'list' },
        })
    );
    
    listUsersTrend.add(listRes.timings.duration);
    
    check(listRes, {
        'list users status is 200': (r) => {
            const isSuccess = r.status === 200;
            isSuccess ? successfulRequests.add(1) : failedRequests.add(1);
            return isSuccess;
        },
    }) || errorRate.add(1);

    // Small sleep to prevent overwhelming the server
    sleep(1);
}

// Handle test summary
export function handleSummary(data) {
    const summary = {
        'create_user': {
            p95: createUserTrend.values.length > 0 ? createUserTrend.p(95).toFixed(2) : 0,
            avg: createUserTrend.values.length > 0 ? createUserTrend.avg.toFixed(2) : 0,
            min: createUserTrend.values.length > 0 ? createUserTrend.min.toFixed(2) : 0,
            max: createUserTrend.values.length > 0 ? createUserTrend.max.toFixed(2) : 0
        },
        'get_user': {
            p95: getUserTrend.values.length > 0 ? getUserTrend.p(95).toFixed(2) : 0,
            avg: getUserTrend.values.length > 0 ? getUserTrend.avg.toFixed(2) : 0,
            min: getUserTrend.values.length > 0 ? getUserTrend.min.toFixed(2) : 0,
            max: getUserTrend.values.length > 0 ? getUserTrend.max.toFixed(2) : 0
        },
        'update_user': {
            p95: updateUserTrend.values.length > 0 ? updateUserTrend.p(95).toFixed(2) : 0,
            avg: updateUserTrend.values.length > 0 ? updateUserTrend.avg.toFixed(2) : 0,
            min: updateUserTrend.values.length > 0 ? updateUserTrend.min.toFixed(2) : 0,
            max: updateUserTrend.values.length > 0 ? updateUserTrend.max.toFixed(2) : 0
        },
        'list_users': {
            p95: listUsersTrend.values.length > 0 ? listUsersTrend.p(95).toFixed(2) : 0,
            avg: listUsersTrend.values.length > 0 ? listUsersTrend.avg.toFixed(2) : 0,
            min: listUsersTrend.values.length > 0 ? listUsersTrend.min.toFixed(2) : 0,
            max: listUsersTrend.values.length > 0 ? listUsersTrend.max.toFixed(2) : 0
        },
        'overall': {
            total_requests: data.metrics.http_reqs.values.count,
            success_rate: (100 - (failedRequests.value / data.metrics.http_reqs.values.count * 100)).toFixed(2) + '%',
            avg_duration: data.metrics.http_req_duration.values.avg.toFixed(2),
            p95_duration: data.metrics.http_req_duration.values.p(95).toFixed(2)
        }
    };

    return {
        'stdout': JSON.stringify(summary, null, 2),
        'summary.json': JSON.stringify(summary, null, 2)
    };
}
