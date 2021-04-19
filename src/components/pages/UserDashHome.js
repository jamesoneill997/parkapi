import React from 'react'
import DashBar from '../DashBar'
import DashBody from '../DashBody'


export default class UserDashHome extends React.Component{
    constructor(props){
        super(props)
    }

    render(){
    return(
        <React.Fragment>
        <DashBar></DashBar>
        <DashBody title={"Welcome to Park.ai"} name={"Paul"}></DashBody>
        </React.Fragment>
    )}
}