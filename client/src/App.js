import React from "react";
import "./App.css";

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      urls: [],
    };
  }

  render() {
    return (
      <div>
        <h1>CargoShortener</h1>
        <NameForm></NameForm>
      </div>
    );
  }
}
class NameForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = { value: "" };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({ value: event.target.value });
  }

  handleSubmit(event) {
    let s = event.target[0].value;
    if (!s.match(/^[a-zA-Z]+:\/\//)) {
      s = "http://" + s;
    }
    fetch("/create/", {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      method: "POST",
      body: JSON.stringify({ longUrl: s }),
    })
      .then((response) => {
        console.log(response);
        if (response.status == 422) {
          alert("Invalid URL!");
        } else {
          response = response.json();
        }
      })
      .then((data) => {
        alert("Your url is now available at: " + data.shortUrl);
      })
      .catch(function (res) {
        console.log(res);
      });
    event.preventDefault();
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit}>
        <label>
          Name:
          <input
            type="text"
            value={this.state.value}
            onChange={this.handleChange}
          />
        </label>
        <input type="submit" value="Submit" />
      </form>
    );
  }
}
export default App;
