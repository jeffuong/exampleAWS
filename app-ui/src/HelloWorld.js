import React from 'react';

/*const HelloWorld = () => {
  
  function sayHello() {
    alert('Hello, World!');
  }
  
  return (
    <button onClick={sayHello}>Click me!</button>
  );
};*/



class HelloWorld extends React.Component {
    constructor() {
        super()
        this.state = {
            s3Objs: []
        }

    }

    componentDidMount() {
        const apiUrl = 'http://localhost:8080/getObjects';
        
        fetch(apiUrl, {mode: 'cors'})
            .then(function(response) {
                return response.json()
            })
            .then((data) => this.setState({ 
                s3Objs: data
            }));
        
            
    }
    

    render() {
        /*return (
            <div>
                <p>Files: {this.state.FileName}</p>           
            </div>
            
        )*/
        return (this.state.s3Objs.map(values => {
            return <div>
                        <p><a href={values.Url}>FileName: {values.FileName}</a></p>
                        <p>StorageType: {values.StorageType} </p>
                        <p>ETag: {values.ETag} </p>
                   </div>
            }

            
            ))
    }

}

export default HelloWorld;