import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Notification } from 'src/app/model/notification';
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

  constructor(private _requestService: RequestService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    this.getNotifications()
  }


  getNotifications() {
    this._requestService.getNotifications(localStorage.getItem("id"))
          .subscribe(data => {
              this.notifications = data.notifications
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
  }

}
