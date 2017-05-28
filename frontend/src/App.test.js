import React from 'react';
import ReactDOM from 'react-dom';
import { expect } from 'chai';
import { spy } from 'sinon';
import { mount } from 'enzyme';
import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';
import App from './App';

// MOCKS

// have to mock local storage because it's not available in jest
class LocalStorageMock {
  constructor() { this.store = {} }

  clear() { this.store = {} }

  getItem(key) { return this.store[key] || null }

  setItem(key, value) { this.store[key] = value }

  removeItem(key) { delete this.store[key] }
}
global.localStorage = new LocalStorageMock

// TESTS

const mock = new MockAdapter(axios);

it('renders without an error', () => {

  mock.onGet('http://localhost:8000/session').reply(200, {
    sessionId: "randomSessionId"
  });
  spy(App.prototype, 'componentDidMount');

  const wrapper = mount(<App />);
  expect(App.prototype.componentDidMount.calledOnce).to.equal(true);
  App.prototype.componentDidMount.restore(); //cleaning the spy
});

// more of a behavioural test, checking that the session is set corectly and then poking the state with a proverbial stick
// to see if everything is set correctly
it('successfully sets the session id', async () => {

  mock.onGet('http://localhost:8000/session').reply(200, {
    sessionId: "randomSessionId"
  });

  spy(App.prototype, 'componentWillMount');
  spy(App.prototype, 'updateSessionId');
  spy(localStorage, 'setItem');

  const wrapper = mount(<App />);
  localStorage.clear();

  // we have no way to wait for the componentWillMount func but we can wait for our helper function
  await wrapper.instance().updateSessionId();
  expect(App.prototype.componentWillMount.calledOnce).to.equal(true);

  // because we manually call it again, there should be 2 calls to update session id
  expect(App.prototype.updateSessionId.callCount).to.equal(2);

  // but because state is changed from the first async call, the second call does not set the session id anymore
  expect(localStorage.setItem.callCount).to.equal(1);
  expect(wrapper.state().sessionId).to.equal('randomSessionId');

  // sinon spy cleanup
  App.prototype.componentWillMount.restore();
  App.prototype.updateSessionId.restore();
  localStorage.setItem.restore();
});

describe('changing inputs', () => {

  it('does allow the values to be changed', async () => {
    const wrapper = mount(<App />);

    const email = wrapper.find('[name="email"]');
    email.simulate('change', { target: { value: 'mail@random.com' }});
    expect(wrapper.state().email).to.equal('mail@random.com');

  });
});
