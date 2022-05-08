import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {

  role: any

  constructor(private router: Router) { }

  ngOnInit(): void {
    this.role = localStorage.getItem('role')
  }

  logOut() {
    localStorage.removeItem('id')
    localStorage.removeItem('jwtToken')
    localStorage.removeItem('role')
    localStorage.removeItem('username')

    this.role = ''

    this.router.navigateByUrl('/');
  }

}
