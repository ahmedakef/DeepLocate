import React, { Component } from 'react'
import axios from 'axios'
import { v4 as uuidv4 } from 'uuid';
import { Button, Row, Container, Form, ButtonToolbar, Col } from 'react-bootstrap';
import DateTimePicker from 'react-datetime-picker';
import './mainForm.css';

class MainForm extends Component {
  state = {
    path: '',
    query: '',
    content: false,
    minSize: 0,
    maxSize: 0,
    minSizeUnit: "Byte",
    maxSizeUnit: "Byte",
    extentions: "",
    minAccessDate: new Date(),
    maxAccessDate: new Date(),
    minModifyDate: new Date(),
    maxModifyDate: new Date(),
    minChangeDate: new Date(),
    maxChangeDate: new Date(),
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

  handleMinSizeChange = (event) => {
    this.setState({ minSize: event.target.value });
  }

  handleMaxSizeChange = (event) => {
    this.setState({ maxSize: event.target.value });
  }

  handleMinSizeUnitChange = (event) => {
    this.setState({ minSizeUnit: event.target.value });
  }

  handleMaxSizeUnitChange = (event) => {
    this.setState({ maxSizeUnit: event.target.value });
  }

  handleExtentionChange = (event) => {
    this.setState({ extentions: event.target.value });
  }

  handleMinAccessDateChange = date => this.setState({ minAccessDate: date })
  handleMaxAccessDateChange = date => this.setState({ maxAccessDate: date })

  handleMinModifyDateChange = date => this.setState({ minModifyDate: date })
  handleMaxModifyDateChange = date => this.setState({ maxModifyDate: date })

  handleMinChangeDateChange = date => this.setState({ minChangeDate: date })
  handleMaxChangeDateChange = date => this.setState({ maxChangeDate: date })

  handleSearch = () => {
    console.log(this.state)

    if (this.state.showAdvanced) {
      //&Atime=2014-11-12T11:45&startSize=2&endSize=800000"
      var url = '/metaSearch?q=' + this.state.query + '&destination=' + this.state.path + '&deepScan=' + this.state.content
      if (this.state.extentions.length > 0) {
        url += '&extentions=' + this.state.extentions
      }
      axios.get(url)
        .then(res => this.setState({ results: [...res.data.matchedFiles] }))
        .catch(err => console.log(err));
    } else {
      axios.get('/search?q=' + this.state.query + '&destination=' + this.state.path + '&deepScan=' + this.state.content)
        .then(res => this.setState({ results: [...res.data.matchedFiles, ...res.data.contentMatchedFiles] }))
        .catch(err => console.log(err));
    }
  }

  handleIndex = () => {
    axios.get('/index?q=' + this.state.query + '&destination=' + this.state.path + '&deepScan=' + this.state.content)
      .then(res => console.log(res))
      .catch(err => console.log(err));
  }

  handleUpdate = () => {
    axios.get('/update?q=' + this.state.query + '&destination=' + this.state.path + '&deepScan=' + this.state.content)
      .then(res => console.log(res))
      .catch(err => console.log(err));
  }

  render() {
    return (
      <Container>
        <Form>
          <Form.Group as={Row} controlId="formPath">
            <Form.Label column sm="2">Path</Form.Label>
            <Col sm="10">
              <Form.Control type="text" value={this.state.path} onChange={this.handlePathChange} placeholder="Enter path to the folder" />
            </Col>
          </Form.Group>

          <Form.Group as={Row} controlId="formQuery">
            <Form.Label column sm="2">Query</Form.Label>
            <Col sm="10">
              <Form.Control type="text" value={this.state.query} onChange={this.handleQueryChange} placeholder="Enter the query to search for" />
            </Col>
          </Form.Group>

          <Form.Check type="checkbox" value={this.state.content} onChange={this.handleContentChange} label="Include Files Contents" />

          <Form.Check type="checkbox" value={this.state.showAdvanced} onChange={this.handleAdvancedChange} label="Advanced" />
        </Form>

        <Col>
          {this.state.showAdvanced ?
            <Container>
              <Form.Group as={Row} controlId="formExtention">
                <Form.Label column sm="2">Extetnions</Form.Label>
                <Col sm="10">
                  <Form.Control type="text" value={this.state.extentions} onChange={this.handleExtentionChange} placeholder="Enter extetntions to filter or leave blank to include all" />
                </Col>
              </Form.Group>
              <Form.Group as={Row} controlId="formSize">
                <Form.Label column sm="1">Size</Form.Label>
                <Col sm="3">
                  <Form.Control type="number" value={this.state.minSize} onChange={this.handleMinSizeChange} />
                </Col>
                <Col sm="2">
                  <Form.Control as="select" onChange={this.handleMinSizeUnitChange} value={this.state.minSizeUnit} custom>
                    <option>Byte</option>
                    <option>KB</option>
                    <option>MB</option>
                    <option>GB</option>
                    <option>TB</option>
                  </Form.Control>
                </Col>
                <Col sm="1">
                  <span> To </span>
                </Col>
                <Col sm="3">
                  <Form.Control type="number" value={this.state.maxSize} onChange={this.handleMaxSizeChange} />
                </Col>
                <Col sm="2">
                  <Form.Control as="select" onChange={this.handleMaxSizeUnitChange} value={this.state.maxSizeUnit} custom>
                    <option>Byte</option>
                    <option>KB</option>
                    <option>MB</option>
                    <option>GB</option>
                    <option>TB</option>
                  </Form.Control>
                </Col>
              </Form.Group>
              <Row>
                <Col>
                  <span>Access Time:</span>
                </Col>
                <Col>
                  <DateTimePicker
                    onChange={this.handleMinAccessDateChange}
                    value={this.state.minAccessDate}
                  />
                </Col>
                <Col>
                  <span> To </span>
                </Col>
                <Col>
                  <DateTimePicker
                    onChange={this.handleMaxAccessDateChange}
                    value={this.state.maxAccessDate}
                  />
                </Col>
              </Row>
              <Row>
                <Col>
                  <span>Modify Time:</span>
                </Col>
                <Col>
                  <DateTimePicker
                    onChange={this.handleMinModifyDateChange}
                    value={this.state.minModifyDate}
                  />
                </Col>
                <Col>
                  <span> To </span>
                </Col>
                <Col>
                  <DateTimePicker
                    onChange={this.handleMaxModifyDateChange}
                    value={this.state.maxModifyDate}
                  />
                </Col>
              </Row>
              <Row>
                <Col>
                  <span>Change Time:</span>
                </Col>
                <Col>
                  <DateTimePicker
                    onChange={this.handleMinChangeDateChange}
                    value={this.state.minChangeDate}
                  />
                </Col>
                <Col>
                  <span> To </span>
                </Col>
                <Col>
                  <DateTimePicker
                    onChange={this.handleMaxChangeDateChange}
                    value={this.state.maxChangeDate}
                  />
                </Col>
              </Row>
            </Container> : null}
        </Col>

        <Row>
          <Col>
            <Button variant="primary" size="lg" onClick={this.handleSearch}>Search</Button>
          </Col>
          <Col>
            <Button variant="primary" size="lg" onClick={this.handleIndex}>Index</Button>
          </Col>
          <Col>
            <Button variant="primary" size="lg" onClick={this.handleUpdate}>Update</Button>
          </Col>
        </Row>

        <div>
          {this.state.results.map(value => <p key={uuidv4()}>{value}</p>)}
        </div>

      </Container >
    );
  }
}


export default MainForm;
