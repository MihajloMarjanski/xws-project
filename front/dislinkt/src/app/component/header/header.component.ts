import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {

  roles: any

  constructor(private router: Router) { }

  ngOnInit(): void {
    this.roles = localStorage.getItem('roles')
  }

  logOut() {
    localStorage.removeItem('id')
    localStorage.removeItem('jwtToken')
    localStorage.setItem('roles', '')
    localStorage.removeItem('authorities')
    localStorage.removeItem('username')

    this.roles = ''

    this.router.navigateByUrl('/');
  }

}
