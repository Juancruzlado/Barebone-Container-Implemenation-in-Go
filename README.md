# Barebone-Container-Implemenation-in-Go

Implement a barebones container.

POC idea: 
Build a basic container runtime with chroot, implementing isolation in networking, security, process management.
Build it all from scratch, in golang. 
Key goal: It should be able to run a service inside an isolated "environment" via tricking the user. Chroot with steroids.

Inspiriations and sources: 

- https://www.youtube.com/watch?v=Utf-A4rODH8  Liz Rice - Building a container from scratch in Go.
- https://www.youtube.com/watch?v=JOsWB50LmwQ Earthly - Build your own Container Runtime.
- https://www.youtube.com/watch?v=sK5i-N34im8 Jérôme Petazzoni - Cgroups, namespaces, and beyond: what are containers made from? 
- https://www.youtube.com/watch?v=GFUpXhft8zA Vishesh @DeepSource - Building containers from scratch 