import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { User } from 'src/app/model/user';
import { RequestService } from 'src/app/service/request.service';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-potential-connections',
  templateUrl: './potential-connections.component.html',
  styleUrls: ['./potential-connections.component.css']
})
export class PotentialConnectionsComponent implements OnInit {

  displayedColumns: string[] = ['name', 'gender', 'email', 'username', 'connect'];
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


  constructor(private _requestService: RequestService, public _userService: UserService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    this.getRecommended()
  }


  getRecommended() {
    this._userService.getRecommendedConnections(localStorage.getItem("id"))
          .subscribe(data => {
              this.users = data.users
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
    
  }

  sendRequest(selectedUser: User) {
    this._requestService.sendConnectRequest(localStorage.getItem("id"), selectedUser.id)
          .subscribe(data => {
            this.getRecommended()
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error); 
  }

}
