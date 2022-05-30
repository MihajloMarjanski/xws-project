import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, FormGroupDirective, NgForm, ValidationErrors, Validators } from '@angular/forms';
import {ErrorStateMatcher} from '@angular/material/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { User } from 'src/app/model/user';
import { UserService } from 'src/app/service/user.service';


/** Error when invalid control is dirty, touched, or submitted. */
export class MyErrorStateMatcher implements ErrorStateMatcher {
  isErrorState(control: FormControl | null, form: FormGroupDirective | NgForm | null): boolean {
    const isSubmitted = form && form.submitted;
    return !!(control && control.invalid && (control.dirty || control.touched || isSubmitted));
  }
}

const specialChars = /[`!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/;


@Component({
  selector: 'app-registration-page',
  templateUrl: './registration-page.component.html',
  styleUrls: ['./registration-page.component.css']
})
export class RegistrationPageComponent implements OnInit {

  emailControl: FormControl = new FormControl('', [Validators.required, Validators.email]);
  //space: FormControl = new FormControl('', [spaceValidator]);

  matcher = new MyErrorStateMatcher();
  user: User = {
    id: 0,
    name: "",
    gender: "Male",
    email: "",
    username:  "",
    password: "",
    biography: '',
    date: new Date(),
    phone: ''
  }
  errorMessage : string  = '';
  repassword: string = '';
  usernames: string[] = [];
  blackListPassMessage: string = '';
  isInBlackList: boolean = false;
  role: string = 'user';

  constructor(public _userService: UserService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    this.getAllUsernames();
  }

  submit(): void {
      this._userService.createUser(this.user)
        .subscribe(
          data => {
              this.router.navigateByUrl('/').then(() => {
                this._snackBar.open('Registration request successfully submited! Activate your account via email.', 'Close', {duration: 5000});
                }); 
            
          },
          error => {
            this._snackBar.open('Invalid input!', 'Close', {duration: 5000});
            console.log('Error!', error)
          }
        )
      
  }

  getAllUsernames() {
    this._userService.getAllUsernames()
        .subscribe(data => {
          this.usernames = data
          console.log(this.usernames);},
                    error => this.errorMessage = <any>error);
  }

  checkPass() {
    this._userService.checkBlackListPass(this.user.password)
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

  containAllCharacters(pass: string) {
    var res = specialChars.test(pass);
    return res
  }

  containSpace(username: string) {
    if(username.split('').includes(' '))
      return true
    return false
  }
}
