import React from 'react'

import '../../App.css'
import HomeSection from '../HomeSection'

export default class Home extends React.Component{
    constructor(props){
        super(props)
    }
    render(){
    return(
        <React.Fragment>
            <HomeSection/>
        </React.Fragment>
    )
    }
}
