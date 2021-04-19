import React from 'react'
import '../App.css'
import {Button} from './Button'
import './HomeSection.css'

function HomeSection() {
    return (
        <div className='home-container'>
        <img src={require("../assets/images/logo.png").default} alt='logo' className='container-logo'></img>
            <h1>Smart parking, made simple.</h1>
            <div className="home-btn">
                <Button className='btns' buttonStyle='btn--outline' buttonSize='btn--large'>GET STARTED</Button>
            </div>
        </div>
    )
}

export default HomeSection
