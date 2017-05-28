import React, { Component } from 'react';
import axios from 'axios';
import ReactResizeDetector from 'react-resize-detector';
import './App.css';

const SESSION_URL = 'http://localhost:8000/session';
const EVENTS_URL = 'http://localhost:8000/events'

class App extends Component {

  constructor(props) {
    super(props);

    this.state = {
      email: '',
      cardNumber: '',
      cvv: '',
      agreedToTerms: false
    };

    this.initialWindowDimensions = {
      height: window.innerHeight,
      width: window.innerWidth
    };

    this.copiedFields = {};

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handlePaste = this.handlePaste.bind(this);
  }

  // we want to be done with the session handling before the component is displayed
  componentWillMount() {

    this.updateSessionId();
  }

  updateSessionId() {
    let sessionId = window.localStorage.getItem('sessionId');

    if(!sessionId) {
      return axios.get(SESSION_URL)
        .then(res => {
          sessionId = res.data.sessionId;
          window.localStorage.setItem('sessionId', sessionId);
          this.setState({ sessionId });
        });
    } else {
      this.setState({ sessionId });
    }
  }

  // because we have autofocus on the email, we can assume the "completionTime" sstarts going as soon
  // as the component is visible
  componentDidMount() {
    this.startTime = Date.now();
  }

  // we use this function to update state, separated reponsibility from listening to paste events
  handleChange(event) {
    const name = event.target.name; // representing the name property on the element that was just changed
    const value = event.target.value; // getting the new value

    this.setState({
      [name]: value
    });
  }

  // this function knows how to handle paste events (and update our backend)
  handlePaste(event) {

    this.copiedFields[event.target.name] = true;

    const payload = {
      sessionId: this.state.sessionId,
      websiteurl: window.location.href,
      eventType: 'copiedFields',
      copiedFields: this.copiedFields
    };

    axios.post(EVENTS_URL, payload)
      .then(response => { console.log(response); })
      .catch(error => { console.error(error); });
  }

  // this function knows how to submit the form 
  // (in this case it actually just sends the time it took to fill the form)
  handleSubmit(event) {
    console.log(event);
    event.preventDefault();

    const finishedTime = Date.now();
    const timeTaken = Math.round( (finishedTime - this.startTime) / 1000 );

    const payload = {
      sessionId: this.state.sessionId,
      websiteUrl: window.location.href,
      eventType: 'submitted',
      timeTaken,
    };

    axios.post(EVENTS_URL, payload)
      .then(response => { console.log(response); })
      .catch(error => { console.error(error); });
  }

  // bit of a hack but we're using this function here because we assume a strict component like behaviour and
  // the final customer might not want to bleed code outside / just want a plug-n-play solution
  // so we have to get the window height ourselves
  onResize() {

    const payload = {
      sessionId: this.state.sessionId,
      websiteUrl: window.location.href,
      eventType: 'resized',
      fromWidth: this.initialWindowDimensions.width,
      fromHeight: this.initialWindowDimensions.height,
      toWidth: window.innerWidth,
      toHeight: window.innerHeight
    };

    axios.post(EVENTS_URL, payload)
      .then(response => { console.log(response); })
      .catch(error => { console.error(error); });
  }

  render() {
    return (
      <div className="App"> 
        <ReactResizeDetector handleWidth handleHeight onResize={this.onResize.bind(this)} />

        <div className="container">
          <pre>Your session id is: {this.state.sessionId}</pre>

          <form className="form-details" onSubmit={this.handleSubmit} method="post">

            <h2 className="form-details-heading">Details</h2>

            <label htmlFor="email" className="sr-only">Email address</label>
            <input type="email"
              name="email"
              className="form-control"
              value={this.state.email}
              onChange={this.handleChange}
              onPaste={this.handlePaste}
              placeholder="Email address"
              required autoFocus />

            <label htmlFor="cardNumber" className="sr-only">Card Number</label>
            <input type="text"
              name="cardNumber"
              value={this.state.cardNumber}
              onChange={this.handleChange}
              onPaste={this.handlePaste}
              className="form-control"
              placeholder="Card Number"
              required />

            <label htmlFor="cvv" className="sr-only">CVV</label>
            <input type="text"
              name="cvv"
              value={this.state.cvv}
              onChange={this.handleChange}
              onPaste={this.handlePaste}
              className="form-control"
              placeholder="Security Code"
              required />

            <div className="checkbox">
              <label>
                <input type="checkbox" value="agree" /> Agree to Terms
              </label>
            </div>

            <button className="btn btn-lg btn-primary btn-block" type="submit">Submit</button>
          </form>

        </div>                     
      </div>
    );
  }
}

export default App;
