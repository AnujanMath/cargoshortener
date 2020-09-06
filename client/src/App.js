import React from "react";
import "./App.css";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import SendIcon from "@material-ui/icons/Send";
import FileCopyIcon from "@material-ui/icons/FileCopy";
import Alert from "@material-ui/lab/Alert";
import Collapse from "@material-ui/core/Collapse";
import Grid from "@material-ui/core/Grid";

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
        <Grid
          container
          width="100%"
          direction="column"
          justify="center"
          alignItems="center"
        >
          <Grid item xs={12}>
            <h1>CargoShortener</h1>
          </Grid>
          <Grid item width="100%" xs={12}>
            <NameForm></NameForm>
          </Grid>
        </Grid>
      </div>
    );
  }
}
class NameForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      value: "",
      button: false,
      error: false,
      collapse: false,
      url: "",
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({ button: false });
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
        if (response.status === 422) {
          this.setState({ error: true });
        } else {
          response.json().then((data) => {
            this.setState({ error: false, url: data.shortUrl });
          });
        }
      })
      .catch(function (res) {
        console.log(res);
      });
    this.setState({ button: true, collapse: true });

    event.preventDefault();
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit} noValidate autoComplete="off">
        <Grid container direction="row" justify="center" alignItems="center">
          {" "}
          <TextField
            error={this.state.error}
            id="outlined-error-helper-text"
            label="Enter URL"
            helperText={this.state.error ? "Invalid URL!" : ""}
            variant="outlined"
            onChange={this.handleChange}
          />
          <Button
            variant="contained"
            disabled={this.state.button}
            color="primary"
            justify="center"
            endIcon={<SendIcon></SendIcon>}
          >
            Send
          </Button>
        </Grid>

        <Collapse in={this.state.collapse}>
          <Alert
            icon={
              <FileCopyIcon
                style={{ cursor: "pointer" }}
                fontSize="inherit"
                onClick={() => {
                  navigator.clipboard.writeText(this.state.url);
                }}
              />
            }
            onClose={() => {
              this.setState({ collapse: false });
            }}
            variant="filled"
            severity="success"
          >
            Your URL is now available at {this.state.url}
          </Alert>
        </Collapse>
      </form>
    );
  }
}
export default App;
