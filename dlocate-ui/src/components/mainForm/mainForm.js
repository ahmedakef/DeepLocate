import React, { Component } from 'react'
import { v4 as uuidv4 } from 'uuid';
import { Button, Row, Container, Form , ButtonToolbar } from 'react-bootstrap';
import './mainForm.css';

class MainForm extends Component {
  state = {
    path: '',
    query: '',
    content: false,
    minSize: -1,
    maxSize: -1,
    extentions: [],
    minAccessDate: null,
    maxAccessDate: null,
    minModifyDate: null,
    maxModifyDate: null,
    minChangeDate: null,
    maxChangeDate: null,
    showAdvanced: false,
    results: [],
  };

  handlePathChange = (event) => {
    this.setState({ path: event.target.value });
  }

  handleQueryChange = (event) => {
    this.setState({ query: event.target.value });
  }

  handleContentChange = (event) => {
    this.setState({ content: event.target.checked });
  }

  handleAdvancedChange = (event) => {
    this.setState({ showAdvanced: event.target.checked });
  }

  handleSearch = () => {
    console.log(this.state)
    this.setState({ results: ["random file", "anything", "just for testing"] });
  }

  handleIndex = () => {
    console.log(this.state)
  }

  handleUpdate = () => {
    console.log(this.state)
  }

  render() {
    return (
      <Container>

        <Form>
          <Form.Group controlId="formPath">
            <Form.Label>Path</Form.Label>
            <Form.Control type="text" value={this.state.path} onChange={this.handlePathChange} placeholder="Enter path to the folder" />
          </Form.Group>

          <Form.Group controlId="formQuery">
            <Form.Label>Query</Form.Label>
            <Form.Control type="text" value={this.state.query} onChange={this.handleQueryChange} placeholder="Enter the query to search for" />
          </Form.Group>

          <Form.Check type="checkbox" value={this.state.content} onChange={this.handleContentChange} label="Include Files Contents" />

          <Form.Check type="checkbox" value={this.state.showAdvanced} onChange={this.handleAdvancedChange} label="Advanced" />
        </Form>

        <Row>
          {this.state.showAdvanced ? <div>test</div> : null}
        </Row>

        <ButtonToolbar>
          <Button variant="primary" onClick={this.handleSearch}>Search</Button>
          <Button variant="primary" onClick={this.handleIndex}>Index</Button>
          <Button variant="primary" onClick={this.handleUpdate}>Update</Button>
        </ButtonToolbar>

        <div>
          {this.state.results.map(value => <p key={uuidv4()}>{value}</p>)}
        </div>

      </Container>
    );
  }
}


export default MainForm;
