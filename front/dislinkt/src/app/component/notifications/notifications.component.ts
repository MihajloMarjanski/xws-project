import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Notification } from 'src/app/model/notification';
import { UsernameWithReq } from 'src/app/model/username-with-req';
import { RequestService } from 'src/app/service/request.service';

@Component({
  selector: 'app-notifications',
  templateUrl: './notifications.component.html',
  styleUrls: ['./notifications.component.css']
})
export class NotificationsComponent implements OnInit {

  displayedColumns: string[] = ['messages', 'date'];
  errorMessage : string  = '';
  notifications : Notification[] = []
  requests : UsernameWithReq[] = []

  constructor(private _requestService: RequestService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    this.getNotifications()
    this.getRequests()
  }


  getNotifications() {
    this._requestService.getNotifications(localStorage.getItem("id"))
          .subscribe(data => {
              this.notifications = data.notifications
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
  }

  getRequests() {
    this._requestService.getRequests(localStorage.getItem("id"))
          .subscribe(data => {
              this.requests = data.users
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
  }

  accept(req: UsernameWithReq) {
    this._requestService.acceptRequest(req.senderId, req.receiverId)
    .subscribe(data => {
        this.getRequests()
        this._snackBar.open('You are connected now :)', 'Close', {duration: 3000});
      console.log('Dobio: ', data)},
    error => this.errorMessage = <any>error);  
  }

  decline(req: UsernameWithReq) {
    this._requestService.declineRequest(req.senderId, req.receiverId)
          .subscribe(data => {
            this.getRequests()
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
  }

}
