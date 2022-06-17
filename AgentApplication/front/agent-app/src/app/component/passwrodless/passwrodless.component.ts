import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { UserService } from 'src/app/service/user.service';
import jwt_decode from 'jwt-decode';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-passwrodless',
  templateUrl: './passwrodless.component.html',
  styleUrls: ['./passwrodless.component.css']
})
export class PasswrodlessComponent implements OnInit {

  constructor(private activatedRoute: ActivatedRoute, private userService: UserService, private router: Router, private _snackBar: MatSnackBar) {
    this.activatedRoute.queryParams.subscribe(params => {
          let token = params['token'];
          console.log(token); // Print the parameter to the console. 

          localStorage.setItem('jwtToken', token)
          userService.loginPasswordless(token).subscribe(
            data => {
              localStorage.setItem('jwtToken', data)
              let tokeninfo = this.getDecodedAccessToken(data)
              if(tokeninfo == '') 
                this._snackBar.open(data, 'Close', {duration: 6000});  
              else {
                localStorage.setItem('username', tokeninfo.sub)
                localStorage.setItem('roles', tokeninfo.roles)
                localStorage.setItem('authorities', tokeninfo.authorities)
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
