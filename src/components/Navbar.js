import React, {useState, useEffect} from 'react'
import './Navbar.css'
import {Button} from './Button.js'
import {Link} from 'react-router-dom'
import '../../node_modules/font-awesome/css/font-awesome.min.css'; 

function Navbar() {
    const [click, setClick] = useState(false)
    const [button, setButton] = useState(true)

    const handleClick = ()=>setClick(!click)
    const closeMobileMenu = ()=>setClick(false)
    const showButton = ()=>{
        if(window.innerWidth<=960){
            setButton(false)
        }else{
            setButton(true)
        }
    }

    useEffect(()=>{
        showButton()
    },[])

    window.addEventListener('resize', showButton)

    return (
        <React.Fragment>
            <nav className="nbar">
                
                <div className="nbar-container">
                    <Link to="/" className="nbar-logo" onClick={closeMobileMenu}>
                    <img src={require("../assets/images/logo.png").default} alt='parkai-logo' />
                    </Link>
                    <div className="menu-icon" onClick={handleClick}>
                        <i className={click ? 'fas fa-times' : 'fas fa-bars'} />
                    </div>
                    <ul className={click?'nav-menu-active':'nav-menu'}>
                        <li className='nav-item'>
                            <Link to='/' className='nav-links' onClick={closeMobileMenu}>
                                Home
                            </Link>
                        </li>
                        <li className='nav-item'>
                            <Link to='/ourjourney' className='nav-links' onClick={closeMobileMenu}>
                                Our Journey
                            </Link>
                        </li>
                        <li className='nav-item'>
                            <Link to='/why' className='nav-links' onClick={closeMobileMenu}>
                                Why Park.AI?
                            </Link>
                        </li>
                        <li className='nav-item'>
                            <Link to='/blog' className='nav-links-mobile' onClick={closeMobileMenu}>
                                Blog
                            </Link>
                        </li>
                    </ul>
                    {button && <Button buttonStyle='btn--outline' buttonSize='btn--medium'>Login/Register</Button>}
                </div>
            </nav>
        </React.Fragment>
    )
}

export default Navbar
