import React from 'react'
import './Form.css'
import Cookies from 'js-cookie'
import jwt from 'jwt-decode'
import {useHistory} from"react-router-dom"


class LoginForm extends React.Component{

    constructor(props) {
        super(props);
        this.firstNameEl = React.createRef();
        this.surnameEl = React.createRef();
        this.emailEl = React.createRef();
        this.passwordEl = React.createRef();
        this.typeEl = React.createRef();
    }

    state = {
        isLoading : false
    }

    handleSubmit = (e) => {

        this.setState({isLoading:true})
        e.preventDefault()
    
        const data = {
            email: this.emailEl.current.value,
            password: this.passwordEl.current.value
        }
        
        fetch('http://localhost:8080/login', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
        
            body: JSON.stringify(data)
        }).then((response)=>{
            
            this.setState({isLoading: false})

            switch (response.status){
                case 200:
                    console.log(response.text().then(function(data){
                        Cookies.set("ParkAIToken", data, {expires: 7, path:'/'})
                        console.log(jwt(Cookies.get('ParkAIToken')))
                    }))


                    fetch('http://localhost:8080/users', {
                        method: 'GET',
                        credentials: 'include',
                        headers: {
                            'Accept': 'application/json',
                            'Content-Type': 'application/json',
                            'Set-Cookie': "ParkAIToken=" + Cookies.get('ParkAIToken')
                        },
                    })

                    this.props.history.push("/dashboard/home")

                    
                    return

                case 401:
                    alert("Invalid credentials")
                    return

                default:
                    alert("Unknown error, please contact support@parkai.com")
                    return
            }
        }).catch((error)=>{
            console.log(error)
        })
        
    }

    render(){

        const {isLoading} = this.state
        return (
            <form onSubmit={this.handleSubmit}>
                    <h3>Sign In</h3>

                    <div className="form-group">
                        <label>Email address</label>
                        <input type="email" className="form-control" placeholder="Enter email" ref={this.emailEl} required/>
                    </div>

                    <div className="form-group">
                        <label>Password</label>
                        <input type="password" className="form-control" placeholder="Enter password" ref={this.passwordEl} required/>
                    </div>

                   
                    <div className="button-container">
                        <button type="submit" className="btn btn-primary btn-block" disabled={isLoading}>
                        {isLoading && <i className="fa fa-refresh fa-spin"></i>}
                        Sign Up
                        </button>
                    </div>
                <div className="form-end">
                    <p className="forgot-password text-right">
                        Don't have an account? <a href="/signup">Sign Up</a>
                        </p>
                    </div>
                </form>
        )}

}


export default LoginForm