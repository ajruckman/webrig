ForEach ($scss in $( Get-ChildItem -Exclude theme.scss | Where-Object { $_.Extension -eq ".scss" } ))
{
    ForEach ($theme in $( Get-ChildItem -Directory .\themes\ ))
    {
        Write-Output $($scss.Name) $($theme.Name)
        S:\libsass\sassc.exe -I . -I .\themes\$($theme.Name)\ .\$( $scss.Name ) .\build\$($scss.BaseName).$( $theme.Name ).css
    }
}
