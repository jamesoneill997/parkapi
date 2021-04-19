import React from 'react'
import DashBar from '../DashBar'
import DashBody from '../DashBody'


export default class UserDashTopup extends React.Component{
    constructor(props){
        super(props)
    }

    render(){
    return(
        <React.Fragment>
        <DashBar></DashBar>
        <DashBody title={"Top Up"} ></DashBody>
        </React.Fragment>
    )}
}