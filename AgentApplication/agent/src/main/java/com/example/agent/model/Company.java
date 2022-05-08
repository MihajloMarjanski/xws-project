package com.example.agent.model;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import javax.persistence.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Entity
public class Company {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;
    private String name;
    private String info;

    @OneToOne(fetch = FetchType.EAGER)
    @JoinColumn(name = "company_owner_id")
    private CompanyOwner companyOwner;

    private boolean isApproved;
}
