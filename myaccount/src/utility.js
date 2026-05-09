import { MILISECONDS_TO_SECONDS } from './constants/constants';
import {
  ENV_NAME_DEVELOPMENT,
  ENV_NAME_STAGING,
  ENV_NAME_PRODUCTION,
  DEVELOPMENT_APEX,
  STAGING_APEX,
  PRODUCTION_APEX,
  HTTP_PREFIX,
  HTTPS_PREFIX,
  getRefetchInterval,
} from './constants/config';

export function emptyPromise(val = null) {
  /* creates an empty promise for cases when data doesn't need to be fetched */
  return new Promise((resolve) => { resolve(val); });
}

export function timestampExpired(timestamp, endpoint = 'DEFAULT') {
  /*
    input: timestamp and a (str) endpoint name
    output: true if the timestamp has elapsed longer than the endpoint allows
  */
  const timeDiff = (Date.now() - timestamp) / MILISECONDS_TO_SECONDS;
  return timeDiff > getRefetchInterval(endpoint);
}

export function detectEnvironment() {
  let env = ENV_NAME_DEVELOPMENT;
  const url = window.location.href.toLowerCase();
  const domain = url.split('/')[2];

  if (domain.endsWith(PRODUCTION_APEX)) {
    env = ENV_NAME_PRODUCTION;
  } else if (domain.endsWith(STAGING_APEX)) {
    env = ENV_NAME_STAGING;
  }

  return env;
}

// devPathPrefix maps a service name to the path Faraday routes to it on
// localhost (path-based routing — see ADR-0004).
const devPathPrefix = {
  whoami: '/whoami',
  account: '/api/account',
  company: '/api/company',
  myaccount: '/myaccount',
  app: '/app',
  ical: '/ical',
  superpowers: '/superpowers',
  www: '',
};

export function routeToMicroservice(service, path = '') {
  switch (detectEnvironment()) {
    case ENV_NAME_DEVELOPMENT: {
      const prefix = devPathPrefix[service] !== undefined
        ? devPathPrefix[service]
        : `/${service}`;
      return prefix + path;
    }

    case ENV_NAME_STAGING:
      return `${HTTPS_PREFIX}${service}${STAGING_APEX}${path}`;

    case ENV_NAME_PRODUCTION:
      return `${HTTPS_PREFIX}${service}${PRODUCTION_APEX}${path}`;

    default: {
      const prefix = devPathPrefix[service] !== undefined
        ? devPathPrefix[service]
        : `/${service}`;
      return prefix + path;
    }
  }
}

export function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response;
  }

  const error = new Error(response.statusText);
  error.response = response;
  throw error;
}

export function parseJSON(response) {
  return response.json();
}
