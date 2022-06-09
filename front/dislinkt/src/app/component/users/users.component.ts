import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { User } from 'src/app/model/user';
import { RequestService } from 'src/app/service/request.service';
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
  isConnected : boolean = false


  constructor(private _requestService: RequestService, public _userService: UserService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    this.searchTable()
    this.role = localStorage.getItem("roles")
  }



  showUser(selectedUser: User) {
    this.user = selectedUser
    this.show = true
    this.areConnected()
  }

  searchTable() {
    this._userService.searchUsers(this.searchField, localStorage.getItem("id"))
          .subscribe(data => {
              this.users = data.users
              this.show = false
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
    
  }


  areConnected() {
    this._requestService.areConnected(localStorage.getItem("id"), this.user.id)
          .subscribe(data => {
            this.isConnected = data.AreConnected
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error); 
  }

  connect() {
    this._requestService.sendConnectRequest(localStorage.getItem("id"), this.user.id)
          .subscribe(data => {
            this.show = false
            this.searchTable()
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error); 
  }

  viewPosts() {
    this.router.navigate(["/posts"],{state: {data: this.user.id}})
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
