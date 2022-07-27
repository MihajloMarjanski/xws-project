import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Client } from 'src/app/model/client';
import { CompanyOwner } from 'src/app/model/company-owner';
import { JobPosition } from 'src/app/model/JobPosition';
import { CompanyService } from 'src/app/service/company.service';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-user-profile',
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.css']
})
export class UserProfileComponent implements OnInit {

  client: Client = {
      id: 0,
      firstName: "",
      lastName: "",
      email: "",
      username:  "",
      password: ""      
  };
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
      ownerUsername: "",
      city: "",
      country: "",
      comments: [],
      positions: []
    }
  }
  repassword: string = "";
  errorMessage : string  = '';
  role: string|null = localStorage.getItem('roles');
  oldPassword:string = ""
  wantToChangePassword : boolean = false;
  newJob: JobPosition = {
    id: 0,
    name: "",
    avgSalary: 0,
    interviewInformations: []
}
selectedPositions : JobPosition[] = []
  

  constructor(public _userService: UserService, private router: Router, _companyService: CompanyService, private _snackBar: MatSnackBar) { }

  ngOnInit(): void {
    this.getUserByUsername();
  }

  addNewJob() {
    let newJob  = Object.assign({},this.newJob);
    this.owner.company.positions.push(newJob);
    this.newJob.name = '';
  }

  changeWantToChangePassword() {
    this.wantToChangePassword = !this.wantToChangePassword
  }

  edit() {
    if(this.role == 'ROLE_CLIENT') 
    {
      if(!this.wantToChangePassword)
        this.client.password = this.oldPassword
      this._userService.editClient(this.client)
          .subscribe(data => {
            console.log('Dobio: ', data)
            if(data == null)
              this._snackBar.open('Incorrect filling of form or someone edited before you! Check and send again edit request', 'Close', {duration: 5000});
            else
              this.client = data
              this.router.navigateByUrl('/')
            },
          error => this.errorMessage = <any>error); 

      console.log(this.client);
      this._snackBar.open('Successfully edited', 'Close', {duration: 5000}); 
    }
    else if(this.role == 'ROLE_ADMIN')
    {
      if(!this.wantToChangePassword)
        this.client.password = this.oldPassword
      this._userService.editAdmin(this.client)
      .subscribe(data => {
        console.log('Dobio: ', data)
        if(data == null)
          this._snackBar.open('Incorrect filling of form or someone edited before you! Check and send again edit request', 'Close', {duration: 5000});
        else
          this.client = data
          this.router.navigateByUrl('/')
        },
      error => this.errorMessage = <any>error); 

      console.log(this.client);
      this._snackBar.open('Successfully edited', 'Close', {duration: 5000}); 
    }
    else if(this.role == 'ROLE_COMPANY_OWNER' || this.role ==='ROLE_POTENTIAL_OWNER')
    {
      this.clientToOwner()
      if(!this.wantToChangePassword)
        this.owner.password = this.oldPassword      
      this._userService.editOwner(this.owner)
          .subscribe(data => {
            console.log('Dobio: ', data)
            if(data == null)
              this._snackBar.open('Incorrect filling of form or someone edited before you! Check and send again edit request', 'Close', {duration: 5000});
            else
              this.client = data
              this.router.navigateByUrl('/')
            },
          error => this.errorMessage = <any>error); 

      console.log(this.client);
      this._snackBar.open('Successfully edited', 'Close', {duration: 5000});  
    }
    
  }


  getUserByUsername()
  {
    if(this.role == 'ROLE_CLIENT') 
    {
      this._userService.getClientByUsername(localStorage.getItem('username') || '')
      .subscribe(data => {
                  this.client = data
                  this.oldPassword = data.password
                  this.client.password = ""
                  this.repassword = ''
                  
                  console.log('Dobio: ', data)},
                error => this.errorMessage = <any>error);  
    }
    else if(this.role == 'ROLE_ADMIN')
    {
      this._userService.getAdminByUsername(localStorage.getItem('username') || '')
      .subscribe(data => {
                  this.client = data
                  this.oldPassword = data.password
                  this.client.password = ""
                  this.repassword = ''
                  console.log('Dobio: ', data)},
                error => this.errorMessage = <any>error);  
    }
    else if(this.role == 'ROLE_COMPANY_OWNER' || this.role ==='ROLE_POTENTIAL_OWNER')
    {
      this._userService.getCompanyOwnerByUsername(localStorage.getItem('username') || '')
      .subscribe(data => {
                  this.owner = data
                  this.oldPassword = data.password
                  this.ownerToClient()
                  this.owner.password = ""
                  this.client.password = ""
                  this.repassword = ''
                  this.selectedPositions = this.owner.company.positions
                  console.log('Dobio: ', data)},
                error => this.errorMessage = <any>error);   
    }    
  }

  ownerToClient()
  {
    this.client.id = this.owner.id;
    this.client.firstName= this.owner.firstName;
    this.client.lastName= this.owner.lastName;
    this.client.email= this.owner.email;
    this.client.username= this.owner.username;
    this.client.password= this.owner.password;
  }

  clientToOwner()
  {
    this.owner.firstName= this.client.firstName;
    this.owner.lastName= this.client.lastName;
    this.owner.password= this.client.password;
  }
}
