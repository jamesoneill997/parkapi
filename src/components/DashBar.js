import React, {useState, useEffect} from 'react'
import {Button} from './Button.js'
import {Link} from 'react-router-dom'
import '../../node_modules/font-awesome/css/font-awesome.min.css'; 
import * as FaIcons from 'react-icons/fa';
import * as AiIcons from 'react-icons/ai';
import { SidebarData } from './SidebarData';
import './DashBar.css';
import { IconContext } from 'react-icons';

function DashBar() {
  return (
    <>
      <IconContext.Provider value={{ color: '#fff' }}>

        <nav className='dash-menu active'>
          <ul className='dash-menu-items'>
            <li className='dashbar-toggle'>
            <img className='dash-logo' src={require("../assets/images/logo.png").default} alt='parkai-logo' />            </li>
            
            {SidebarData.map((item, index) => {
              return (
                <li key={index} className={item.cName}>
                  <Link to={item.path}>
                    {item.icon}
                    <span>{item.title}</span>
                  </Link>
                </li>
              );
            })}
          </ul>
        </nav>
      </IconContext.Provider>
    </>
  );
}

export default DashBar;