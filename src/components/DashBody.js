import React, {useState, useEffect} from 'react'
import {Button} from './Button.js'
import {Link} from 'react-router-dom'
import '../../node_modules/font-awesome/css/font-awesome.min.css'; 
import * as FaIcons from 'react-icons/fa';
import * as AiIcons from 'react-icons/ai';
import { SidebarData } from './SidebarData';
import './DashBody.css';
import { IconContext } from 'react-icons';

export default class DashBody extends React.Component{
  render(){

    return (
    <div className="body-container">
        <h1>{this.props.title} {this.props.name}</h1>
    
    </div>
    
  )}
}

