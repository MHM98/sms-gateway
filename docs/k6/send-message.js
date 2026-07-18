import http from 'k6/http';
import { check } from 'k6';

const baseURL = __ENV.BASE_URL || 'http://127.0.0.1:3000';
const duration = __ENV.DURATION || '1m';
const numberOfUsers = 1000;

function createUserIDs(number) {
  return Array.from({ length: number }, (_, index) => index + 1);
}

const userIDs = createUserIDs(numberOfUsers);

export const options = {
  discardResponseBodies: true,

  scenarios: {
    normal_messages: {
      executor: 'constant-arrival-rate',
      exec: 'sendNormal',
      rate: 1500,
      timeUnit: '1s',
      duration,
      preAllocatedVUs: 450,
      maxVUs: 1500,
      tags: {
        service_type: 'normal',
      },
    },

    express_messages: {
      executor: 'constant-arrival-rate',
      exec: 'sendExpress',
      rate: 500,
      timeUnit: '1s',
      duration,
      preAllocatedVUs: 150,
      maxVUs: 500,
      tags: {
        service_type: 'express',
      },
    },
  },

  thresholds: {
    'http_req_failed{service_type:normal}': ['rate<0.01'],
    'http_req_failed{service_type:express}': ['rate<0.01'],

    'http_req_duration{service_type:normal}': [
      'p(95)<200',
      'p(99)<500',
    ],
    'http_req_duration{service_type:express}': [
      'p(95)<200',
      'p(99)<500',
    ],

    dropped_iterations: ['count<1'],
  },
};

function sendMessage(serviceType) {
  const userID = userIDs[Math.floor(Math.random() * userIDs.length)];

  const response = http.post(
    `${baseURL}/api/v1/sms-gateway/message/send`,
    JSON.stringify({
      user_id: userID,
      recipient: `+98912${String(Math.floor(Math.random() * 10000000)).padStart(7, '0')}`,
      body: `k6 ${serviceType} load-test message`,
      service_type: serviceType,
    }),
    {
      headers: {
        'Content-Type': 'application/json',
        'X-Idempotency-Key': crypto.randomUUID(),
      },
      tags: {
        service_type: serviceType,
        name: 'POST /api/v1/sms-gateway/message/send',
      },
      timeout: '5s',
    },
  );

  check(response, {
    'status is 204': (result) => result.status === 204,
  });
}

export function sendNormal() {
  sendMessage('normal');
}

export function sendExpress() {
  sendMessage('express');
}
