import React from 'react'
import './App.css';
import Navbar from './components/Navbar'
import {BrowserRouter as Router, Switch, Route} from 'react-router-dom'
import './App.css'
import Home from './components/pages/Home'
import Signup from './components/pages/Signup'
import SignIn from './components/pages/SignIn'
import UserDashHome from './components/pages/UserDashHome'
import UserDashParking from './components/pages/UserDashParking'
import UserDashGarage from './components/pages/UserDashGarage'
import UserDashAccount from './components/pages/UserDashAccount'
import UserDashTopup from './components/pages/UserDashTopup'


export default class App extends React.Component {
    constructor(){
        super()
        this.state = {
            loggedIn:false,
            user: {}
        }
    }

    render(){
  return (
      <React.Fragment>
        <Router>
            <Switch>
            <Route path='/dashboard/home' 
            exact 
            render={props=>(
                <UserDashHome {... props} loggedIn = {this.state.loggedIn} />
            )}/>
            <Route path='/dashboard/parking' 
            exact 
            render={props=>(
                <UserDashParking {... props} loggedIn = {this.state.loggedIn} />
            )}/>
            <Route path='/dashboard/garage' 
            exact 
            render={props=>(
                <UserDashGarage {... props} loggedIn = {this.state.loggedIn} />
            )}/>
            <Route path='/dashboard/topup' 
            exact 
            render={props=>(
                <UserDashTopup {... props} loggedIn = {this.state.loggedIn} />
            )}/>
            <Route path='/dashboard/account' 
            exact 
            render={props=>(
                <UserDashAccount {... props} loggedIn = {this.state.loggedIn} />
            )}/>

            <div>
            <Navbar/>

                <Route 
                path='/' 
                exact 
                render={props=>(
                    <Home {... props} loggedIn = {this.state.loggedIn} />
                )}
                />
                <Route path='/signup'             
                exact 
                render={props=>(
                    <Signup {... props} loggedIn = {this.state.loggedIn} />
                )}/>
                <Route path='/signin' 
                exact 
                render={props=>(
                    <SignIn {... props} loggedIn = {this.state.loggedIn} />
                )}/>
                </div>
            </Switch>
        </Router>
      </React.Fragment>
  )}
}

