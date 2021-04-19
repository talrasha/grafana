declare let __webpack_public_path__: string;

/**
 * Check if we are hosting files on cdn and set webpack public path
 */
if ((window as any).public_cdn_path) {
  __webpack_public_path__ = (window as any).public_cdn_path;
}

import 'symbol-observable';
import 'core-js';
import 'regenerator-runtime/runtime';

import 'whatwg-fetch'; // fetch polyfill needed for PhantomJs rendering

import _ from 'lodash';
import ReactDOM from 'react-dom';
import React from 'react';
import AppWrapper from './AppWrapper';

ReactDOM.render(React.createElement(AppWrapper), document.getElementById('reactRoot'));
