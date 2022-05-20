package com.example.agent.mapper;

import com.example.agent.model.Company;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.dto.OwnerWithCompany;

import java.util.Date;

public class CompanyOwnerAdapter {

    public static OwnerWithCompany convertOwnerToDto(CompanyOwner owner, Company company) {
        OwnerWithCompany dto = new OwnerWithCompany();

        dto.setId(owner.getId());
        dto.setFirstName(owner.getFirstName());
        dto.setLastName(owner.getLastName());
        dto.setEmail(owner.getEmail());
        dto.setUsername(owner.getUsername());
        dto.setPassword(owner.getPassword());
        dto.setActivated(owner.isActivated());
        dto.setForgotten(owner.getForgotten());
        dto.setPin(owner.getPin());
        dto.setSalt(owner.getSalt());
        dto.setMissedPasswordCounter(owner.getMissedPasswordCounter());
        dto.setBlocked(owner.isBlocked());
        dto.setBlockedDate(owner.getBlockedDate());

        dto.setCompany(company);

        return dto;
    }
}
