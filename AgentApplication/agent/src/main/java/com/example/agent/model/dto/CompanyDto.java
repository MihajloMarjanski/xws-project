package com.example.agent.model.dto;

import com.example.agent.model.Comment;
import com.example.agent.model.Company;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.JobPosition;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import javax.persistence.*;
import java.util.HashSet;
import java.util.Set;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class CompanyDto {
    private Integer id;
    private String name;
    private String info;
    private String city;
    private String country;
    private CompanyOwner companyOwner;
    private boolean isApproved;

    public CompanyDto(Company company) {
        id = company.getId();
        name = company.getName();
        info = company.getInfo();
        city = company.getCity();
        country = company.getCountry();
        companyOwner = company.getCompanyOwner();
        isApproved = company.isApproved();
    }

}
