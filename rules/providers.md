# Providers

- A provider is an Interface the core declares  
  and something outside fulfils.
- The core Depends on the shape.  
  Never on the library behind it.
- Name the Provider after the business need,  
  never after the vendor.  
  StudentStore, not GormRepository.
- One struct may Fulfil several providers.  
  One provider must never Leak its vendor.
- Comment each Provider with the URL  
  of the contract it wraps.  
  A reader should not have to Search.
- Count the Files that import a vendor.  
  If the count grows past one, the provider Failed.
- The word Collides with Angular, NestJS and Terraform,  
  where a provider is a registered dependency.  
  Here it is a Port.
