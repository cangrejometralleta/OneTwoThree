# Providers

- A Provider is an Interface the Core Declares  
  and something outside Fulfils.
- The Core Depends on the Shape.  
  Never on the Library behind it.
- Name the Provider after the Business Need,  
  never after the Vendor.  
  StudentStore, not GormRepository.
- One Struct May Fulfil several Providers.  
  One Provider Must never Leak its Vendor.
- Comment each Provider with the URL  
  of the Contract it Wraps.  
  A Reader Should not have to Search.
- Count the Files that Import a Vendor.  
  If the Count Grows past one, the Provider Failed.
- The Word Collides with Angular, NestJS and Terraform,  
  where a Provider is a registered Dependency.  
  Here it is a Port.
