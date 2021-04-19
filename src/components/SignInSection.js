import React from 'react'
import LoginForm from './LoginForm'
import '../App.css'
import './HomeSection.css'
import './LoginForm.css'
import {useHistory} from "react-router-dom";

function SignInSection(){

    return (
        <div className='home-container'>
            <LoginForm history={useHistory()}></LoginForm>
        </div>
    )
}

export default SignInSection