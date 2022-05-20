import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Client } from 'src/app/model/client';
import { CompanyOwner } from 'src/app/model/company-owner';
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
      ownerUsername: ""
    }
  }
  repassword: string = "";
  errorMessage : string  = '';
  role: string|null = localStorage.getItem('role');

  constructor(public _userService: UserService, _companyService: CompanyService, private _snackBar: MatSnackBar) { }

  ngOnInit(): void {
    this.getUserByUsername();
  }


  edit() {
    if(this.role == 'ROLE_CUSTOMER') 
    {
      this._customerService.edit(this.customer)
          .subscribe(data => {
            console.log('Dobio: ', data)
            if(data == null)
              this._snackBar.open('Incorrect filling of form or someone edited before you! Check and send again edit request', 'Close', {duration: 5000});
            else
              this.customer = data
            },
          error => this.errorMessage = <any>error); 

      console.log(this.customer);
      this._snackBar.open('Successfully edited', 'Close', {duration: 5000}); 
    }
    else if(this.role == 'ROLE_ADMIN')
    {

    }
    else if(this.role == 'ROLE_BOAT_OWNER')
    {
      //this.customerToNonCustomer();
      this._boatOwnerService.edit(this.nonCustomer)
          .subscribe(data => {
            console.log('Dobio: ', data)
            if(data == null)
              this._snackBar.open('Incorrect filling of form! Check and send again edit request', 'Close', {duration: 5000});
            else
            {
              this.nonCustomer = data;
              this.nonCustomerToCustomer();//sad su u customeru sva bitna polja za prikaz od nonCustomera
            }
            },
          error => this.errorMessage = <any>error); 

      console.log(this.customer);
      this._snackBar.open('Successfully edited', 'Close', {duration: 5000});   
    }
    else if(this.role == 'ROLE_WEEKEND_HOUSE_OWNER')
    {
      //this.customerToNonCustomer();
      this._weekendHouseOwnerService.edit(this.nonCustomer)
          .subscribe(data => {
            console.log('Dobio: ', data)
            if(data == null)
              this._snackBar.open('Incorrect filling of form! Check and send again edit request', 'Close', {duration: 5000});
            else
            {
              this.nonCustomer = data;
              this.repassword = this.nonCustomer.password;
              this.nonCustomerToCustomer();//sad su u customeru sva bitna polja za prikaz od nonCustomera
            }
            },
          error => this.errorMessage = <any>error); 

      console.log(this.customer);
      this._snackBar.open('Successfully edited', 'Close', {duration: 5000});   
    }
    else if(this.role == 'ROLE_INSTRUCTOR')
    {
      
    }  
  }


  getUserByUsername()
  {
    if(this.role == 'ROLE_CLIENT') 
    {
      this._userService.getClientByUsername(localStorage.getItem('username') || '')
      .subscribe(data => {
                  this.client = data
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
                  this.client.password = ""
                  this.repassword = ''
                  
                  console.log('Dobio: ', data)},
                error => this.errorMessage = <any>error);  
    }
    else if(this.role == 'ROLE_COMPANY_OWNER')
    {
      this._userService.getCompanyOwnerByUsername(localStorage.getItem('username') || '')
      .subscribe(data => {
                  this.owner = data
                  this.owner.password = ""
                  this.repassword = ''
                  console.log('Dobio: ', data)},
                error => this.errorMessage = <any>error);   
    }    
  }
}
