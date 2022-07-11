import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PotentialConnectionsComponent } from './potential-connections.component';

describe('PotentialConnectionsComponent', () => {
  let component: PotentialConnectionsComponent;
  let fixture: ComponentFixture<PotentialConnectionsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PotentialConnectionsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PotentialConnectionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
