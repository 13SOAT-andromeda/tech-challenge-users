import http from 'k6/http';
import { check, sleep } from 'k6';

// ===============================
// CONFIGURAÇÕES
// ===============================

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8081';
const USER_ID = __ENV.USER_ID || '1';
const USER_EMAIL = __ENV.USER_EMAIL || 'admin@admin.com';
const USER_ROLE = __ENV.USER_ROLE || 'administrator';

// Parâmetros de configuração de thresholds via variáveis de ambiente
const MAX_VUS = Number(__ENV.MAX_VUS || 100);  // número máximo de usuários virtuais
const P95_MS = Number(__ENV.P95_MS || 800);    // as requisições devem ser respondidas em até 800ms
const ERROR_RATE = Number(__ENV.ERROR_RATE || 0.01); // taxa máxima de erro permitida (1%)

export const options = {
  stages: [
    { duration: '30s', target: Math.floor(MAX_VUS * 0.2) },
    { duration: '1m',  target: Math.floor(MAX_VUS * 0.5) },
    { duration: '1m',  target: MAX_VUS },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: [`p(95)<${P95_MS}`],
    http_req_failed: [`rate<${ERROR_RATE}`],
  },
};

export default function () {
  // O serviço tech-challenge-users valida identidade via headers injetados pelo Lambda Authorizer.
  // Em testes de stress, os headers são fornecidos diretamente.
  const headers = {
    'X-User-Id': USER_ID,
    'X-User-Email': USER_EMAIL,
    'X-User-Role': USER_ROLE,
    'Content-Type': 'application/json',
  };

  // -------- USERS --------
  const usersRes = http.get(
    `${BASE_URL}/users`,
    { headers }
  );

  check(usersRes, {
    'users status 200': (r) => r.status === 200,
  });

  // -------- CUSTOMERS --------
  const customersRes = http.get(
    `${BASE_URL}/users/customers`,
    { headers }
  );

  check(customersRes, {
    'customers status 200': (r) => r.status === 200,
  });

  // -------- VEHICLES --------
  const vehiclesRes = http.get(
    `${BASE_URL}/users/vehicles`,
    { headers }
  );

  check(vehiclesRes, {
    'vehicles status 200': (r) => r.status === 200,
  });

  // -------- EMPLOYEES --------
  const employeeRes = http.get(
    `${BASE_URL}/users/employees/1`,
    { headers }
  );

  check(employeeRes, {
    'employee getById status 200 or 404': (r) => r.status === 200 || r.status === 404,
  });

  sleep(1);
}
