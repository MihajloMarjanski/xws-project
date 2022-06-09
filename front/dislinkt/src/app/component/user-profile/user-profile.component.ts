import { Expression } from '@angular/compiler';
import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Experience } from 'src/app/model/experience';
import { Interest } from 'src/app/model/interest';
import { User } from 'src/app/model/user';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-user-profile',
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.css']
})
export class UserProfileComponent implements OnInit {

  user: User = {
    id: 0,
    name: "",
    gender: "Male",
    email: "",
    username:  "",
    password: "",
    biography: '',
    date: new Date(),
    phone: '',
    experiences: [],
    interests: [],
    isPrivate: false
  };
  repassword: string = "";
  errorMessage : string  = '';
  role: string|null = localStorage.getItem('roles');
  oldPassword:string = ""
  wantToChangePassword : boolean = false;
  experiences: Experience[] = []
  interests: Interest[] = []
  displayedColumns: string[] = ['company', 'position', 'from', 'until', 'remove'];
  newExperience: Experience = {
    company: '',
    from: new Date(),
    until: new Date(),
    id: 0,
    position: '',
    userId: 0
  }
  newInterest: Interest = {
    id: 0,
    interest: '',
    userId: 0
  }

  constructor(public _userService: UserService, private router: Router, private _snackBar: MatSnackBar) { }

  ngOnInit(): void {
    this.getUserByUsername();
  }


  changePublic() {
    this.user.isPrivate = !this.user.isPrivate
  }

  addInterest() {
    this.newInterest.userId = this.user.id
    this._userService.addInterest(this.newInterest)
        .subscribe(data => {
                this.getUserByUsername()
                this.newInterest.interest = ''
                this._snackBar.open('Successfully added', 'Close', {duration: 3000});
              },
          error => this.errorMessage = <any>error);
  }

  removeInterest(interest : Interest) {
    this._userService.removeInterest(interest.id)
        .subscribe(data => {
                this._snackBar.open('Successfully removed', 'Close', {duration: 3000});
                this.getUserByUsername()   
              },
          error => this.errorMessage = <any>error);    
  }

  addExperience() {
    this.newExperience.userId = this.user.id
    this._userService.addExperience(this.newExperience)
        .subscribe(data => {
              this.getUserByUsername()
               this.newExperience.company = ''
               this.newExperience.from = new Date(),
               this.newExperience.until = new Date(),
               this.newExperience.position = ''
               this._snackBar.open('Successfully added', 'Close', {duration: 3000});
              },
          error => this.errorMessage = <any>error);
  }

  changeWantToChangePassword() {
    this.wantToChangePassword = !this.wantToChangePassword
  }

  removeExperience(experience : Experience) {
    this._userService.removeExperience(experience.id)
        .subscribe(data => {
                this._snackBar.open('Successfully removed', 'Close', {duration: 3000});
                this.getUserByUsername()    
              },
          error => this.errorMessage = <any>error);    
  }

  edit() {
      if(!this.wantToChangePassword)
        this.user.password = this.oldPassword
      this._userService.editUser(this.user)
          .subscribe(data => {
            console.log('Dobio: ', data)
            if(data.id == 0)
              this._snackBar.open('Incorrect filling of form or someone edited before you! Check and send again edit request', 'Close', {duration: 5000});
            else
              this.user = data
              this.router.navigateByUrl('/').then(() => {
                this._snackBar.open('Successfully edited', 'Close', {duration: 5000}); 
              })
            },
          error => this.errorMessage = <any>error); 

      console.log(this.user);
      
  }
    

  getUserByUsername()
  {
      this._userService.getUserByUsername(localStorage.getItem('username') || '')
      .subscribe(data => {
                  this.user = data.user
                  this.user.date = new Date()
                  this.oldPassword = data.user.password
                  this.user.password = ""
                  this.repassword = ''
                  this.experiences = data.user.experience
                  this.interests = data.user.interests
                  
                  console.log('Dobio: ', data.user)},
                error => this.errorMessage = <any>error);  
    
  }

}
