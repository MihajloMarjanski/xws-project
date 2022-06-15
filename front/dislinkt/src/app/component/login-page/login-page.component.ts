import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Credentials } from '../../model/credentials';
import jwt_decode from 'jwt-decode';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.css']
})
export class LoginPageComponent implements OnInit {

  credentials: Credentials = {
    username: '',
    password: '',
    pin: ''
  }

  token: string = '';
  errorMessage : string  = '';

  constructor(public _userService: UserService, private router: Router, private _snackBar: MatSnackBar) { }

  ngOnInit(): void {
  }


  submit(): void {
    this._userService.logIn(this.credentials).subscribe(      
      data => {
        localStorage.setItem('jwtToken', data.token)
        let tokeninfo = this.getDecodedAccessToken(data.token)
        if(tokeninfo == '') 
          this._snackBar.open("Invalid username or password", 'Close', {duration: 6000});  
        else {
          localStorage.setItem('username', tokeninfo.username)
          localStorage.setItem('roles', tokeninfo.role)
          localStorage.setItem('authorities', tokeninfo.authorities)
          localStorage.setItem('id', tokeninfo.id)
          localStorage.setItem('exp', tokeninfo.exp)
          console.log('Dobio: ', data)
          this.router.navigateByUrl('/').then(() => {
            window.location.reload();
          });
        }
      },
      error => {
        console.log('Error!', error)
        this._snackBar.open(error.error.message, 'Close', {duration: 3000});
      }
    )
  }

  getDecodedAccessToken(token: string): any {
    try{
        return jwt_decode(token);
    }
    catch(Error){
        return '';
    }
  }

  sendPassword() {
    this._userService.forgotPassword(this.credentials.username).subscribe(
      data => {
        this._snackBar.open('Your new password has been sent to your email. You have to change it when you login for the first time!', 'Close', {duration: 7000});
      },
      error => {
        console.log('Error!', error)
        this.errorMessage = <any>error
        this._snackBar.open(error.error.message, 'Close', {duration: 6000});
      }
      )
  }

}
