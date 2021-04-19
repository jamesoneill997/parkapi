import React from 'react'
import DashBar from '../DashBar'
import DashBody from '../DashBody'


export default class UserDashGarage extends React.Component{
    constructor(props){
        super(props)
    }

    render(){
    return(
        <React.Fragment>
        <DashBar></DashBar>
        <DashBody title={"Garage"} ></DashBody>
        </React.Fragment>
    )}
}