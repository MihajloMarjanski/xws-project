import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { Message } from 'src/app/model/message';
import { User } from 'src/app/model/user';
import { RequestService } from 'src/app/service/request.service';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.css']
})
export class ChatComponent implements OnInit {

  selectedUser: User = {
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
  connections : User[] = []
  messages : Message[] = []
  newMessage: Message = {
    receiverId : 0,
    senderId : 0,
    text : ''
  }
  logedUserId : number = 0
  errorMessage : string  = '';

  constructor(private _requestService: RequestService, public _userService: UserService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    this.logedUserId = Number(localStorage.getItem("id"))
    this.findConnections()
  }


  findConnections() {
    this._requestService.findConnections(this.logedUserId)
              .subscribe(data => {
                this.connections = data.users
                console.log('Dobio: ', data)},
              error => this.errorMessage = <any>error);  
  }

  getMessages() {
    this._requestService.getMessagesForConnection(this.logedUserId, this.selectedUser.id)
              .subscribe(data => {
                this.messages = data.messages
                console.log('Dobio: ', data)},
              error => this.errorMessage = <any>error);
  }

  sendMessage() {
    this.newMessage.senderId = this.logedUserId
    this.newMessage.receiverId = this.selectedUser.id
    this._requestService.sendMessage(this.newMessage)
    .subscribe(data => {
            this.newMessage.text = ''
            this.getMessages()
          },
    error => this.errorMessage = <any>error);
  }

}
