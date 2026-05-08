import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import * as Sentry from '@sentry/react';
import configureStore from './stores/configureStore';
import App from './components/App';
import { detectEnvironment } from './utility';
import {
  ENV_NAME_DEVELOPMENT,
  ENV_NAME_PRODUCTION,
  SENTRY_STAGING_KEY,
  SENTRY_PRODUCTION_KEY,
} from './constants/config';

require('./main.scss');

const currentEnv = detectEnvironment();

if (currentEnv !== ENV_NAME_DEVELOPMENT) {
  const sentryDSN = (currentEnv === ENV_NAME_PRODUCTION)
    ? SENTRY_PRODUCTION_KEY : SENTRY_STAGING_KEY;

  Sentry.init({
    dsn: sentryDSN,
    environment: currentEnv,
  });
}

const store = configureStore();

ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('app'),
);
