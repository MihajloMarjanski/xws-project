import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { User } from 'src/app/model/user';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css']
})
export class UsersComponent implements OnInit {

  searchField: string = '';
  displayedColumns: string[] = ['name', 'gender', 'email', 'username'];
  errorMessage : string  = '';
  show : boolean = false
  users: User[] = []
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
  }
  role : string|null = ''

  constructor(public _userService: UserService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    this.searchTable()
    this.role = localStorage.getItem("roles")
  }



  showUser(selectedUser: User) {
    this.user = selectedUser
    this.show = true
  }

  searchTable() {
    this._userService.searchUsers(this.searchField, localStorage.getItem("id"))
          .subscribe(data => {
              this.users = data.users
              this.filterForBlocked()
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
    
  }

  filterForBlocked() {
    this._userService.findAllBlocked(localStorage.getItem("id"))
          .subscribe(data => {
            this.users = this.users.filter( ( el ) => !data.users.includes( el ) );
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
  }

  isConnected() : boolean {
    return false
  }

  connect() {

  }

  block() {
    this._userService.blockUser(localStorage.getItem("id"), this.user.id)
    .subscribe(data => {
        this.show = false
        this.searchTable()
      console.log('Dobio: ', data)},
    error => this.errorMessage = <any>error);  
  }
}
