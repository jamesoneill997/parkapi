import React from 'react'
import DashBar from '../DashBar'
import DashBody from '../DashBody'


export default class UserDashAccount extends React.Component{
    constructor(props){
        super(props)
    }

    render(){
    return(
        <React.Fragment>
        <DashBar></DashBar>
        <DashBody title={"Manage Account"} ></DashBody>
        </React.Fragment>
    )}
}