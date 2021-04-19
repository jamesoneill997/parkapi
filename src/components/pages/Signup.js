import React from 'react'

import '../../App.css'
import SignupSection from '../SignupSection'

export default class Signup extends React.Component{
    constructor(props){
        super(props)
    }
    render(){
    return(
        <React.Fragment>
            <SignupSection/>
        </React.Fragment>
    )
    }
}