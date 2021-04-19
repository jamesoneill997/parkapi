import React from 'react'

import '../../App.css'
import SignInSection from '../SignInSection'

export default class SignIn extends React.Component{
    constructor(props){
        super(props)
        this.handleLogin = this.handleLogin.bind(this)
    }

    handleLogin(data){
        this.props.history.push("/dashboard")
    }

    render(){
    return(
        <React.Fragment>
            <SignInSection handleLogin = {this.handleLogin}/>
        </React.Fragment>
    )}
}