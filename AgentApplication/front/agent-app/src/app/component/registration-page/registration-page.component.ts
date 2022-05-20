import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, FormGroupDirective, NgForm, ValidationErrors, Validators } from '@angular/forms';
import {ErrorStateMatcher} from '@angular/material/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Client } from 'src/app/model/client';
import { CompanyOwner } from 'src/app/model/company-owner';
import { UserService } from 'src/app/service/user.service';


/** Error when invalid control is dirty, touched, or submitted. */
export class MyErrorStateMatcher implements ErrorStateMatcher {
  isErrorState(control: FormControl | null, form: FormGroupDirective | NgForm | null): boolean {
    const isSubmitted = form && form.submitted;
    return !!(control && control.invalid && (control.dirty || control.touched || isSubmitted));
  }
}
// export class UsernameValidator {
//   static cannotContainSpace(control: AbstractControl) : { [key: string]: boolean } | null {
//       if((control.value as string).indexOf(' ') >= 0){
//           return {cannotContainSpace: true}
//       }

//       return null
//   }
// }

// function spaceValidator(control: AbstractControl): { [key: string]: boolean } | null {
//   if((control.value as string).indexOf(' ') >= 0){
//               return {cannotContainSpace: true}
//           }
    
//           return null
// }

@Component({
  selector: 'app-registration-page',
  templateUrl: './registration-page.component.html',
  styleUrls: ['./registration-page.component.css']
})
export class RegistrationPageComponent implements OnInit {

  emailControl: FormControl = new FormControl('', [Validators.required, Validators.email]);
  //space: FormControl = new FormControl('', [spaceValidator]);

  matcher = new MyErrorStateMatcher();
  client: Client = {
    id: 0,
    firstName: "",
    lastName: "",
    email: "",
    username:  "",
    password: "",
  }
  errorMessage : string  = '';
  repassword: string = '';
  usernames: string[] = [];
  blackListPassMessage: string = '';
  isInBlackList: boolean = false;
  role: string = 'client';
  owner: CompanyOwner = {
    id: 0,
    firstName: "",
    lastName: "",
    email: "",
    username: "",
    password: "",
    company: {
      id: 0,
      name: "",
      info: "",
      isApproved: true,
      ownerUsername: ""
    }
  }

  constructor(public _userService: UserService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    this.getAllUsernames();
  }

  submit(): void {
    if (this.role == 'owner') 
    {
      this.owner.email = this.client.email
      this.owner.firstName = this.client.firstName
      this.owner.lastName = this.client.lastName
      this.owner.username = this.client.username
      this.owner.password = this.client.password
      this._userService.createCompanyOwner(this.owner)
        .subscribe(
          data => {
            if(data != null)
              this._snackBar.open(data, 'Close', {duration: 5000});
            else {
              this.router.navigateByUrl('/').then(() => {
                this._snackBar.open('Registration request successfully submited!', 'Close', {duration: 5000});
                }); 
            }
          },
          error => {
            this._snackBar.open('Invalid input! Email already exists.', 'Close', {duration: 5000});
            console.log('Error!', error)
          }
        )
    }  
    else
    {
      this._userService.createClient(this.client)
        .subscribe(
          data => {
            if(data != null)
              this._snackBar.open(data, 'Close', {duration: 5000});
            else {
              this.router.navigateByUrl('/').then(() => {
                this._snackBar.open('Registration request successfully submited!', 'Close', {duration: 5000});
                }); 
            }
          },
          error => {
            this._snackBar.open('Invalid input! Email already exists.', 'Close', {duration: 5000});
            console.log('Error!', error)
          }
        )
    }  
  }

  getAllUsernames() {
    this._userService.getAllUsernames()
        .subscribe(data => {
          this.usernames = data
          console.log(this.usernames);},
                    error => this.errorMessage = <any>error);
  }

  checkPass() {
    this._userService.checkBlackListPass(this.client.password)
        .subscribe(data => {
          if (data == null)
            this.isInBlackList = false
          else {
            this.isInBlackList = true
            this.blackListPassMessage = data
          }
          console.log(this.blackListPassMessage);},
                    error => this.errorMessage = <any>error);
  }


  containSpace(username: string) {
    if(username.split('').includes(' '))
      return true
    return false
  }
}
