import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { UserService } from 'src/app/service/user.service';
import jwt_decode from 'jwt-decode';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-passwordless',
  templateUrl: './passwordless.component.html',
  styleUrls: ['./passwordless.component.css']
})
export class PasswordlessComponent implements OnInit {

  constructor(private activatedRoute: ActivatedRoute, private userService: UserService, private router: Router, private _snackBar: MatSnackBar) {
    this.activatedRoute.queryParams.subscribe(params => {
          let token = params['token'];
          console.log(token); // Print the parameter to the console. 

          localStorage.setItem('jwtToken', token)
          userService.loginPasswordless(token).subscribe(
            data => {
              localStorage.setItem('jwtToken', data.jwt)
              let tokeninfo = this.getDecodedAccessToken(data.jwt)
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
              this._snackBar.open('Invalid token.', 'Close', {duration: 3000});
            }
            )
      });
  }

  ngOnInit(): void {    
  }

  getDecodedAccessToken(token: string): any {
    try{
        return jwt_decode(token);
    }
    catch(Error){
        return '';
    }
  }

}
